'use client'

import { useState, useEffect, useRef, useCallback } from 'react'
import Nav from '@/components/Nav'
import styles from './garage.module.css'

// ── Types ─────────────────────────────────────────────────────────────────────

interface VehicleSpec {
  make: string
  model: string
  year: number
  trim: string
  engine: string
  displacement: number
  cylinders: number
  hp: number
  torque: number
  topSpeed: number
  weight: number
  zeroToSixty: number
  drivetrain: string
  fuelType: string
  seats: number
}

interface TunedStats {
  hp: number
  torque: number
  topSpeed: number
  zeroToSixty: number
  weight: number
  config: string
}

interface TuneResponse {
  base: VehicleSpec
  tuned: TunedStats
  delta: { hp: number; torque: number; topSpeed: number; zeroToSixty: number; weight: number }
  config: TuningConfig
}

interface TuningConfig {
  label: string
  hpMult: number
  torqueMult: number
  topSpeedMult: number
  zeroMult: number
  weightMult: number
}

interface CarEntry {
  make: string
  model: string
  year: number
  category: string
  icon: string
  accentColor: string
  // populated after API fetch
  spec?: VehicleSpec
}

// ── Constants ─────────────────────────────────────────────────────────────────

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

// Hero cars shown in the sidebar. Specs are lazy-fetched from the Go backend.
const HERO_CARS: CarEntry[] = [
  { make: 'Ferrari',     model: 'F40',         year: 1992, category: 'SUPERCAR', icon: '🏎️', accentColor: '#ff2a2a' },
  { make: 'Porsche',     model: '911 GT3 RS',  year: 2024, category: 'TRACK',    icon: '🔶', accentColor: '#fb923c' },
  { make: 'Nissan',      model: 'GT-R',        year: 2023, category: 'SPORTS',   icon: '⚡', accentColor: '#00e5ff' },
  { make: 'Lamborghini', model: 'Huracan Evo', year: 2023, category: 'SUPERCAR', icon: '🐂', accentColor: '#facc15' },
  { make: 'Toyota',      model: 'GR Supra',    year: 2024, category: 'SPORTS',   icon: '🌀', accentColor: '#60a5fa' },
  { make: 'Lotus',       model: 'Evora GT',    year: 2022, category: 'TRACK',    icon: '🍃', accentColor: '#a78bfa' },
]

const TUNING_KEYS = ['stock', 'street', 'track', 'race', 'drift'] as const
type TuningKey = typeof TUNING_KEYS[number]

const TUNING_LABELS: Record<TuningKey, string> = {
  stock: 'STOCK', street: 'STREET', track: 'TRACK', race: 'RACE SPEC', drift: 'DRIFT',
}

// Client-side multipliers — mirrors Go backend, used for instant UI response
// while the real tune POST resolves in the background.
const LOCAL_TUNING: Record<TuningKey, TuningConfig> = {
  stock:  { label: 'Stock',     hpMult: 1.00, torqueMult: 1.00, topSpeedMult: 1.00, zeroMult: 1.00, weightMult: 1.00 },
  street: { label: 'Street',    hpMult: 1.08, torqueMult: 1.06, topSpeedMult: 1.04, zeroMult: 0.95, weightMult: 0.97 },
  track:  { label: 'Track',     hpMult: 1.18, torqueMult: 1.12, topSpeedMult: 1.10, zeroMult: 0.86, weightMult: 0.90 },
  race:   { label: 'Race Spec', hpMult: 1.35, torqueMult: 1.25, topSpeedMult: 1.18, zeroMult: 0.76, weightMult: 0.82 },
  drift:  { label: 'Drift',     hpMult: 1.20, torqueMult: 1.30, topSpeedMult: 0.96, zeroMult: 0.92, weightMult: 0.94 },
}

// ── API helpers ───────────────────────────────────────────────────────────────

async function fetchVehicle(make: string, model: string, year: number): Promise<VehicleSpec | null> {
  try {
    const res = await fetch(
      `${API_BASE}/api/garage/vehicle?make=${encodeURIComponent(make)}&model=${encodeURIComponent(model)}&year=${year}`
    )
    if (!res.ok) return null
    const data = await res.json()
    return data.vehicle ?? null
  } catch {
    return null
  }
}

async function postTune(make: string, model: string, year: number, tuning: TuningKey): Promise<TuneResponse | null> {
  try {
    const res = await fetch(`${API_BASE}/api/garage/tune`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ make, model, year, tuning }),
    })
    if (!res.ok) return null
    return await res.json()
  } catch {
    return null
  }
}

// ── Helpers ───────────────────────────────────────────────────────────────────

function localTune(spec: VehicleSpec, key: TuningKey) {
  const cfg = LOCAL_TUNING[key]
  return {
    hp:          Math.round(spec.hp * cfg.hpMult),
    torque:      Math.round(spec.torque * cfg.torqueMult),
    topSpeed:    Math.round(spec.topSpeed * cfg.topSpeedMult),
    zeroToSixty: +(spec.zeroToSixty * cfg.zeroMult).toFixed(1),
    weight:      Math.round(spec.weight * cfg.weightMult),
  }
}

// ── Sub-components ────────────────────────────────────────────────────────────

function StatBar({ label, value, max, accent }: { label: string; value: number; max: number; accent: string }) {
  const pct = Math.min((value / max) * 100, 100)
  return (
    <div className={styles.statRow}>
      <span className={styles.statKey}>{label}</span>
      <div className={styles.barWrap}>
        <div
          className={styles.barFill}
          style={{ width: `${pct}%`, background: accent, boxShadow: `0 0 8px ${accent}55` }}
        />
      </div>
      <span className={styles.statNum}>{value.toLocaleString()}</span>
    </div>
  )
}

// ── Main page ─────────────────────────────────────────────────────────────────

export default function GaragePage() {
  const [cars, setCars]           = useState<CarEntry[]>(HERO_CARS)
  const [selected, setSelected]   = useState<CarEntry>(HERO_CARS[0])
  const [spec, setSpec]           = useState<VehicleSpec | null>(null)
  const [tuning, setTuning]       = useState<TuningKey>('stock')
  const [tuneResult, setTuneResult] = useState<TuneResponse | null>(null)
  const [loadingSpec, setLoadingSpec] = useState(false)
  const [search, setSearch]       = useState('')
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const scanRef   = useRef<number>(0)
  const rafRef    = useRef<number>(0)

  const accent = selected.accentColor

  // ── Fetch spec when car changes ──────────────────────────────────────────
  useEffect(() => {
    let cancelled = false
    setSpec(null)
    setTuneResult(null)
    setLoadingSpec(true)

    fetchVehicle(selected.make, selected.model, selected.year).then(v => {
      if (cancelled) return
      setSpec(v)
      setLoadingSpec(false)
    })

    return () => { cancelled = true }
  }, [selected])

  // ── Re-tune when tuning or spec changes ─────────────────────────────────
  useEffect(() => {
    if (!spec || tuning === 'stock') { setTuneResult(null); return }
    let cancelled = false

    postTune(selected.make, selected.model, selected.year, tuning).then(r => {
      if (!cancelled) setTuneResult(r)
    })

    return () => { cancelled = true }
  }, [spec, tuning, selected])

  // ── Canvas grid ──────────────────────────────────────────────────────────
  const drawGrid = useCallback(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const W = canvas.offsetWidth
    const H = canvas.offsetHeight
    if (!W || !H) return
    canvas.width  = W
    canvas.height = H

    ctx.clearRect(0, 0, W, H)
    const g = 38

    ctx.strokeStyle = accent + '18'
    ctx.lineWidth   = 0.5
    for (let x = 0; x <= W; x += g) { ctx.beginPath(); ctx.moveTo(x, 0); ctx.lineTo(x, H); ctx.stroke() }
    for (let y = 0; y <= H; y += g) { ctx.beginPath(); ctx.moveTo(0, y); ctx.lineTo(W, y); ctx.stroke() }

    const hor = H * 0.52
    const vpx = W / 2

    ctx.strokeStyle = accent + '30'
    ctx.lineWidth   = 0.8
    for (let i = -14; i <= 14; i++) {
      const x = vpx + i * 58
      ctx.beginPath(); ctx.moveTo(vpx, hor); ctx.lineTo(x, H); ctx.stroke()
    }
    for (let i = 0; i <= 10; i++) {
      const t = i / 10
      const y = hor + (H - hor) * t
      const s = t * W * 0.9
      ctx.beginPath(); ctx.moveTo(vpx - s / 2, y); ctx.lineTo(vpx + s / 2, y); ctx.stroke()
    }

    const gr = ctx.createRadialGradient(vpx, H * 0.6, 0, vpx, H * 0.6, W * 0.5)
    gr.addColorStop(0, accent + '1a')
    gr.addColorStop(1, 'transparent')
    ctx.fillStyle = gr
    ctx.fillRect(0, 0, W, H)
  }, [accent])

  useEffect(() => { drawGrid() }, [drawGrid])
  useEffect(() => {
    const ro = new ResizeObserver(drawGrid)
    if (canvasRef.current) ro.observe(canvasRef.current.parentElement!)
    return () => ro.disconnect()
  }, [drawGrid])

  // ── Scan line animation ──────────────────────────────────────────────────
  useEffect(() => {
    const tick = () => {
      scanRef.current = (scanRef.current + 0.35) % 100
      const el = document.getElementById('scan-line')
      if (el) el.style.top = scanRef.current + '%'
      rafRef.current = requestAnimationFrame(tick)
    }
    rafRef.current = requestAnimationFrame(tick)
    return () => cancelAnimationFrame(rafRef.current)
  }, [])

  // ── Derived display values ───────────────────────────────────────────────
  const displaySpec = spec ?? {
    hp: 0, torque: 0, topSpeed: 0, zeroToSixty: 0, weight: 0,
    engine: '—', drivetrain: '—', displacement: 0,
  }

  const tuned = spec
    ? (tuneResult
        ? tuneResult.tuned
        : localTune(spec as VehicleSpec, tuning))
    : null

  const activeHP    = tuned?.hp    ?? displaySpec.hp
  const activeTorq  = tuned?.torque ?? displaySpec.torque
  const activeSpd   = tuned?.topSpeed ?? displaySpec.topSpeed
  const activeZero  = tuned?.zeroToSixty ?? displaySpec.zeroToSixty
  const activeWt    = tuned?.weight  ?? displaySpec.weight

  const delta = tuneResult?.delta ?? null

  const filtered = search
    ? cars.filter(c =>
        c.make.toLowerCase().includes(search.toLowerCase()) ||
        c.model.toLowerCase().includes(search.toLowerCase()))
    : cars

  return (
    <>
      <Nav />
      <main className={styles.page}>

        {/* ── Sub-header ── */}
        <div className={styles.subHeader}>
          <div className={styles.subHeaderLeft}>
            <span className={styles.breadcrumb}>/ GARAGE</span>
            <span className={styles.subHeaderCar} style={{ color: accent }}>
              {selected.make.toUpperCase()} {selected.model.toUpperCase()}
            </span>
            {loadingSpec && <span className={styles.loadingPip} style={{ background: accent }} />}
          </div>
          <div className={styles.headerStats}>
            <div className={styles.hStat}>
              <span className={styles.hStatN}>{HERO_CARS.length}</span>
              <span className={styles.hStatL}>VEHICLES</span>
            </div>
            <div className={styles.hStatDiv} />
            <div className={styles.hStat}>
              <span className={styles.hStatN}>{TUNING_KEYS.length}</span>
              <span className={styles.hStatL}>CONFIGS</span>
            </div>
            <div className={styles.hStatDiv} />
            <div className={styles.hStat}>
              <span className={styles.hStatN} style={{ color: accent }}>{activeHP || '—'}</span>
              <span className={styles.hStatL}>ACTIVE HP</span>
            </div>
          </div>
        </div>

        <div className={styles.garageGrid}>

          {/* ── LEFT: Car selector ── */}
          <div className={styles.selectorPanel}>
            <div className={styles.panelLabel}>SELECT VEHICLE</div>
            <div className={styles.searchWrap}>
              <input
                className={styles.searchInput}
                type="text"
                placeholder="Search make or model…"
                value={search}
                onChange={e => setSearch(e.target.value)}
              />
            </div>
            <div className={styles.carList}>
              {filtered.map(car => {
                const isActive = car === selected
                const a = car.accentColor
                return (
                  <div
                    key={`${car.make}-${car.model}-${car.year}`}
                    className={`${styles.carSlot} ${isActive ? styles.carSlotActive : ''}`}
                    style={isActive ? { borderColor: `${a}60`, background: `${a}0d` } : {}}
                    onClick={() => { setSelected(car); setTuning('stock') }}
                  >
                    {isActive && (
                      <div className={styles.carSlotAccentBar} style={{ background: a, boxShadow: `0 0 10px ${a}` }} />
                    )}
                    <span className={styles.carSlotIcon}>{car.icon}</span>
                    <div className={styles.carSlotInfo}>
                      <div className={styles.carSlotMake}>{car.make}</div>
                      <div className={styles.carSlotModel}>{car.model}</div>
                      <div className={styles.carSlotYear}>{car.year}</div>
                    </div>
                    <div
                      className={styles.catBadge}
                      style={{ color: a, borderColor: `${a}44`, background: `${a}12` }}
                    >
                      {car.category}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>

          {/* ── CENTER: Viewport ── */}
          <div className={styles.viewportPanel}>
            <div className={styles.viewport}>
              <canvas ref={canvasRef} className={styles.gridCanvas} />
              <div id="scan-line" className={styles.scanLine} style={{ background: `linear-gradient(90deg,transparent,${accent}50,transparent)` }} />

              {/* HUD corners */}
              {(['TL','TR','BL','BR'] as const).map(pos => (
                <div key={pos} className={`${styles.hudCorner} ${styles[`hudCorner${pos}` as keyof typeof styles]}`}>
                  <svg viewBox="0 0 22 22" fill="none" width="22" height="22">
                    <path
                      d={pos.startsWith('T') ? 'M1 21V4L4 1H21' : 'M1 1V18L4 21H21'}
                      stroke={accent} strokeWidth="1.5"
                    />
                  </svg>
                </div>
              ))}

              {/* HUD overlays */}
              <div className={styles.hudTL} style={{ color: accent }}>
                <div className={styles.hudLine}><span className={styles.hudK}>MODEL</span><span>{selected.make} {selected.model}</span></div>
                <div className={styles.hudLine}><span className={styles.hudK}>YEAR</span><span>{selected.year}</span></div>
                <div className={styles.hudLine}><span className={styles.hudK}>ENGINE</span><span>{displaySpec.engine}</span></div>
              </div>
              <div className={styles.hudTR} style={{ color: accent }}>
                <div className={styles.hudLine}><span>{displaySpec.drivetrain}</span><span className={styles.hudK}>DRIVE</span></div>
                <div className={styles.hudLine}><span>{displaySpec.displacement ? `${displaySpec.displacement}cc` : '—'}</span><span className={styles.hudK}>CC</span></div>
                <div className={styles.hudLine}><span>{TUNING_LABELS[tuning]}</span><span className={styles.hudK}>CONFIG</span></div>
              </div>

              {/* Car display */}
              <div className={styles.carDisplay}>
                <div className={styles.glowRing} style={{ boxShadow: `0 0 80px ${accent}44, 0 0 160px ${accent}18` }} />
                <div className={styles.carEmoji} style={{ filter: `drop-shadow(0 0 30px ${accent})` }}>
                  {selected.icon}
                </div>
                <div
                  className={styles.carReflection}
                  style={{ background: `radial-gradient(ellipse 60% 20% at 50% 100%, ${accent}28, transparent)` }}
                />
              </div>

              {/* Status bar */}
              <div className={styles.hudBottom}>
                <span className={styles.statusDot} style={{ background: accent, boxShadow: `0 0 6px ${accent}` }} />
                <span style={{ color: accent }}>
                  {loadingSpec ? 'FETCHING TELEMETRY…' : `SYSTEM ONLINE — ${TUNING_LABELS[tuning]}`}
                </span>
              </div>
            </div>

            {/* Instruments */}
            <div className={styles.instruments} style={{ borderColor: `${accent}25` }}>
              {[
                { val: activeSpd,  label: 'TOP SPEED MPH' },
                { val: activeHP,   label: 'HORSEPOWER' },
                { val: activeZero ? `${activeZero}s` : '—', label: '0 – 60 MPH' },
                { val: activeWt ? activeWt.toLocaleString() : '—', label: 'WEIGHT LBS' },
              ].map(({ val, label }) => (
                <div key={label} className={styles.instrItem}>
                  <span className={styles.instrVal} style={{ color: accent }}>{val || '—'}</span>
                  <span className={styles.instrLabel}>{label}</span>
                </div>
              ))}
            </div>

            {/* Tuning selector */}
            <div className={styles.tuningBar}>
              <span className={styles.tuningLabel}>CONFIGURATION</span>
              <div className={styles.tuningBtns}>
                {TUNING_KEYS.map(key => (
                  <button
                    key={key}
                    className={`${styles.tuningBtn} ${tuning === key ? styles.tuningBtnActive : ''}`}
                    style={tuning === key ? { borderColor: accent, color: accent, background: `${accent}15` } : {}}
                    onClick={() => setTuning(key)}
                  >
                    {TUNING_LABELS[key]}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* ── RIGHT: Stats ── */}
          <div className={styles.statsPanel}>
            <div className={styles.panelLabel}>PERFORMANCE DATA</div>

            <div className={styles.statGroup}>
              <div className={styles.statGroupLabel}>POWER</div>
              <StatBar label="HP"     value={activeHP}   max={800} accent={accent} />
              <StatBar label="TORQUE" value={activeTorq} max={600} accent={accent} />
            </div>

            <div className={styles.statGroup}>
              <div className={styles.statGroupLabel}>DYNAMICS</div>
              <StatBar label="TOP SPD" value={activeSpd} max={250}  accent={accent} />
              <StatBar label="WEIGHT"  value={activeWt}  max={5000} accent={accent} />
            </div>

            <div className={styles.statGroup}>
              <div className={styles.statGroupLabel}>ACCELERATION</div>
              <div className={styles.bigStatCard} style={{ borderColor: `${accent}30`, background: `${accent}08` }}>
                <span className={styles.bigStatNum} style={{ color: accent }}>{activeZero || '—'}</span>
                <span className={styles.bigStatUnit}>sec 0 → 60</span>
              </div>
            </div>

            {/* Delta panel — shown when tuned and backend returned exact delta */}
            {tuning !== 'stock' && spec && (
              <div className={styles.deltaPanel}>
                <div className={styles.deltaPanelLabel}>
                  DELTA VS STOCK
                  {!tuneResult && <span className={styles.deltaEstimate}> (est.)</span>}
                </div>
                <div className={styles.deltaRow}>
                  <span>POWER</span>
                  <span className={styles.pos}>+{activeHP - spec.hp} hp</span>
                </div>
                <div className={styles.deltaRow}>
                  <span>WEIGHT</span>
                  <span className={activeWt < spec.weight ? styles.pos : styles.neg}>
                    {activeWt < spec.weight ? '−' : '+'}{Math.abs(activeWt - spec.weight).toLocaleString()} lbs
                  </span>
                </div>
                <div className={styles.deltaRow}>
                  <span>0–60</span>
                  <span className={styles.pos}>
                    −{(spec.zeroToSixty - activeZero).toFixed(1)}s
                  </span>
                </div>
                {delta && (
                  <div className={styles.deltaRow}>
                    <span>TOP SPEED</span>
                    <span className={styles.pos}>+{delta.topSpeed} mph</span>
                  </div>
                )}
              </div>
            )}

            <div className={styles.apiNote}>
              Data via <a href="https://www.carqueryapi.com" target="_blank" rel="noreferrer">CarQuery API</a>
            </div>
          </div>

        </div>
      </main>
    </>
  )
}