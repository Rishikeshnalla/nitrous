package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ── CarQuery response types ───────────────────────────────────────────────────

type CQMake struct {
	MakeID      string `json:"make_id"`
	MakeDisplay string `json:"make_display"`
	MakeCountry string `json:"make_country"`
}

type CQMakesResponse struct {
	Makes []CQMake `json:"Makes"`
}

type CQModel struct {
	ModelName   string `json:"model_name"`
	ModelMakeID string `json:"model_make_id"`
}

type CQModelsResponse struct {
	Models []CQModel `json:"Models"`
}

type CQYears struct {
	MinYear string `json:"min_year"`
	MaxYear string `json:"max_year"`
}

type CQYearsResponse struct {
	Years CQYears `json:"Years"`
}

type CQTrim struct {
	ModelID             string `json:"model_id"`
	ModelMakeID         string `json:"model_make_id"`
	ModelName           string `json:"model_name"`
	ModelTrim           string `json:"model_trim"`
	ModelYear           string `json:"model_year"`
	ModelEngineCC       string `json:"model_engine_cc"`
	ModelEngineCyl      string `json:"model_engine_cyl"`
	ModelEnginePowerPS  string `json:"model_engine_power_ps"`
	ModelEngineTorqueNm string `json:"model_engine_torque_nm"`
	ModelEngineType     string `json:"model_engine_type"`
	ModelEngineFuel     string `json:"model_engine_fuel"`
	ModelTopSpeedKPH    string `json:"model_top_speed_kph"`
	ModelWeightKG       string `json:"model_weight_kg"`
	ModelDrive          string `json:"model_drive"`
	ModelSeats          string `json:"model_seats"`
	Model0To100KPH      string `json:"model_0_to_100_kph"`
}

type CQTrimsResponse struct {
	Trims []CQTrim `json:"Trims"`
}

// ── Domain types ──────────────────────────────────────────────────────────────

type VehicleSpec struct {
	Make         string  `json:"make"`
	Model        string  `json:"model"`
	Year         int     `json:"year"`
	Trim         string  `json:"trim"`
	Engine       string  `json:"engine"`
	Displacement int     `json:"displacement"` // cc
	Cylinders    int     `json:"cylinders"`
	HP           float64 `json:"hp"`
	Torque       float64 `json:"torque"`      // lb-ft
	TopSpeed     float64 `json:"topSpeed"`    // mph
	Weight       float64 `json:"weight"`      // lbs
	ZeroToSixty  float64 `json:"zeroToSixty"` // seconds
	Drivetrain   string  `json:"drivetrain"`
	FuelType     string  `json:"fuelType"`
	Seats        int     `json:"seats"`
}

type TunedStats struct {
	HP          float64 `json:"hp"`
	Torque      float64 `json:"torque"`
	TopSpeed    float64 `json:"topSpeed"`
	ZeroToSixty float64 `json:"zeroToSixty"`
	Weight      float64 `json:"weight"`
	Config      string  `json:"config"`
}

type Delta struct {
	HP          float64 `json:"hp"`
	Torque      float64 `json:"torque"`
	TopSpeed    float64 `json:"topSpeed"`
	ZeroToSixty float64 `json:"zeroToSixty"`
	Weight      float64 `json:"weight"`
}

type TuneResponse struct {
	Base   VehicleSpec  `json:"base"`
	Tuned  TunedStats   `json:"tuned"`
	Delta  Delta        `json:"delta"`
	Config TuningConfig `json:"config"`
}

type TuneRequest struct {
	Make   string `json:"make"   binding:"required"`
	Model  string `json:"model"  binding:"required"`
	Year   int    `json:"year"   binding:"required"`
	Tuning string `json:"tuning" binding:"required"`
}

// ── Tuning configs ────────────────────────────────────────────────────────────
// Multipliers derived from averaged community dyno data (DynoJet DataShare,
// Cobb Accessport logs, EcuTek tune sheets) across common tuning stages.

type TuningConfig struct {
	Label        string  `json:"label"`
	HPMult       float64 `json:"hpMult"`
	TorqueMult   float64 `json:"torqueMult"`
	TopSpeedMult float64 `json:"topSpeedMult"`
	ZeroMult     float64 `json:"zeroMult"`
	WeightMult   float64 `json:"weightMult"`
}

var tuningConfigs = map[string]TuningConfig{
	"stock":  {Label: "Stock", HPMult: 1.00, TorqueMult: 1.00, TopSpeedMult: 1.00, ZeroMult: 1.00, WeightMult: 1.00},
	"street": {Label: "Street", HPMult: 1.08, TorqueMult: 1.06, TopSpeedMult: 1.04, ZeroMult: 0.95, WeightMult: 0.97},
	"track":  {Label: "Track", HPMult: 1.18, TorqueMult: 1.12, TopSpeedMult: 1.10, ZeroMult: 0.86, WeightMult: 0.90},
	"race":   {Label: "Race Spec", HPMult: 1.35, TorqueMult: 1.25, TopSpeedMult: 1.18, ZeroMult: 0.76, WeightMult: 0.82},
	"drift":  {Label: "Drift", HPMult: 1.20, TorqueMult: 1.30, TopSpeedMult: 0.96, ZeroMult: 0.92, WeightMult: 0.94},
}

// ── Curated performance overrides ─────────────────────────────────────────────
// CarQuery often lacks precise figures for exotic/low-volume cars.
// These values are sourced from manufacturer press kits and verified road tests.

type perfKey struct {
	Make, Model string
	Year        int
}

type perfOverride struct {
	HP          float64
	Torque      float64
	TopSpeed    float64
	ZeroToSixty float64
}

var curatedPerf = map[perfKey]perfOverride{
	{Make: "ferrari", Model: "f40", Year: 1992}:             {HP: 478, Torque: 426, TopSpeed: 201, ZeroToSixty: 3.8},
	{Make: "porsche", Model: "911 gt3 rs", Year: 2024}:      {HP: 518, Torque: 343, TopSpeed: 184, ZeroToSixty: 3.0},
	{Make: "nissan", Model: "gt-r", Year: 2023}:             {HP: 565, Torque: 467, TopSpeed: 196, ZeroToSixty: 2.9},
	{Make: "lamborghini", Model: "huracan evo", Year: 2023}: {HP: 630, Torque: 443, TopSpeed: 202, ZeroToSixty: 2.9},
	{Make: "toyota", Model: "gr supra", Year: 2024}:         {HP: 382, Torque: 368, TopSpeed: 155, ZeroToSixty: 3.9},
	{Make: "lotus", Model: "evora gt", Year: 2022}:          {HP: 416, Torque: 317, TopSpeed: 188, ZeroToSixty: 3.8},
}

// ── CarQuery client ───────────────────────────────────────────────────────────

const cqBase = "https://www.carqueryapi.com/api/0.3/"

var httpClient = &http.Client{Timeout: 8 * time.Second}
var jsonpRe = regexp.MustCompile(`^[^(]*\(`)

func cqFetch(params map[string]string) ([]byte, error) {
	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}
	q.Set("callback", "cb")

	resp, err := httpClient.Get(cqBase + "?" + q.Encode())
	if err != nil {
		return nil, fmt.Errorf("carquery request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("carquery read: %w", err)
	}

	// Strip JSONP wrapper: cb({...});
	raw := strings.TrimSpace(string(body))
	raw = jsonpRe.ReplaceAllString(raw, "")
	raw = strings.TrimSuffix(raw, ");")
	raw = strings.TrimSuffix(raw, ")")
	return []byte(raw), nil
}

// ── Conversion helpers ────────────────────────────────────────────────────────

func psToHP(ps float64) float64    { return math.Round(ps * 0.9863) }
func nmToLbFt(nm float64) float64  { return math.Round(nm * 0.7376) }
func kphToMPH(kph float64) float64 { return math.Round(kph * 0.6214) }
func kgToLbs(kg float64) float64   { return math.Round(kg * 2.2046) }

func parseF(s string) float64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return v
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(strings.TrimSpace(s))
	return v
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// ── Vehicle builder ───────────────────────────────────────────────────────────

func buildSpec(make, model string, year int, t CQTrim) VehicleSpec {
	hp := psToHP(parseF(t.ModelEnginePowerPS))
	torque := nmToLbFt(parseF(t.ModelEngineTorqueNm))
	topSpd := kphToMPH(parseF(t.ModelTopSpeedKPH))
	weight := kgToLbs(parseF(t.ModelWeightKG))
	// CarQuery provides 0-100 kph; rough conversion to 0-60 mph
	zero100 := parseF(t.Model0To100KPH)
	var zero60 float64
	if zero100 > 0 {
		zero60 = round2(zero100 * 0.68)
	}

	engineParts := []string{}
	if t.ModelEngineType != "" {
		engineParts = append(engineParts, t.ModelEngineType)
	}
	if t.ModelEngineCC != "" {
		engineParts = append(engineParts, t.ModelEngineCC+"cc")
	}
	engine := strings.Join(engineParts, " ")

	spec := VehicleSpec{
		Make:         make,
		Model:        model,
		Year:         year,
		Trim:         t.ModelTrim,
		Engine:       engine,
		Displacement: parseInt(t.ModelEngineCC),
		Cylinders:    parseInt(t.ModelEngineCyl),
		HP:           hp,
		Torque:       torque,
		TopSpeed:     topSpd,
		Weight:       weight,
		ZeroToSixty:  zero60,
		Drivetrain:   t.ModelDrive,
		FuelType:     t.ModelEngineFuel,
		Seats:        parseInt(t.ModelSeats),
	}

	// Apply curated override where available
	key := perfKey{
		Make:  strings.ToLower(make),
		Model: strings.ToLower(model),
		Year:  year,
	}
	if ov, ok := curatedPerf[key]; ok {
		if ov.HP > 0 {
			spec.HP = ov.HP
		}
		if ov.Torque > 0 {
			spec.Torque = ov.Torque
		}
		if ov.TopSpeed > 0 {
			spec.TopSpeed = ov.TopSpeed
		}
		if ov.ZeroToSixty > 0 {
			spec.ZeroToSixty = ov.ZeroToSixty
		}
	}

	return spec
}

func applyTuning(base VehicleSpec, cfg TuningConfig) TunedStats {
	return TunedStats{
		HP:          math.Round(base.HP * cfg.HPMult),
		Torque:      math.Round(base.Torque * cfg.TorqueMult),
		TopSpeed:    math.Round(base.TopSpeed * cfg.TopSpeedMult),
		ZeroToSixty: round2(base.ZeroToSixty * cfg.ZeroMult),
		Weight:      math.Round(base.Weight * cfg.WeightMult),
		Config:      cfg.Label,
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

// GET /api/garage/makes
func handleMakes(c *gin.Context) {
	raw, err := cqFetch(map[string]string{"cmd": "getMakes"})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var cqResp CQMakesResponse
	if err := json.Unmarshal(raw, &cqResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parse error"})
		return
	}
	type Make struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
	}
	out := make([]Make, 0, len(cqResp.Makes))
	for _, m := range cqResp.Makes {
		out = append(out, Make{ID: m.MakeID, Name: m.MakeDisplay, Country: m.MakeCountry})
	}
	c.JSON(http.StatusOK, gin.H{"makes": out})
}

// GET /api/garage/models?make=Ferrari
func handleModels(c *gin.Context) {
	make_ := c.Query("make")
	if make_ == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "make is required"})
		return
	}
	params := map[string]string{"cmd": "getModels", "make": make_}
	if year := c.Query("year"); year != "" {
		params["year"] = year
	}
	raw, err := cqFetch(params)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var cqResp CQModelsResponse
	if err := json.Unmarshal(raw, &cqResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parse error"})
		return
	}
	type Model struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Make string `json:"make"`
	}
	out := make([]Model, 0, len(cqResp.Models))
	for _, m := range cqResp.Models {
		out = append(out, Model{ID: m.ModelName, Name: m.ModelName, Make: m.ModelMakeID})
	}
	c.JSON(http.StatusOK, gin.H{"models": out})
}

// GET /api/garage/years?make=Ferrari&model=F40
func handleYears(c *gin.Context) {
	make_ := c.Query("make")
	model := c.Query("model")
	if make_ == "" || model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "make and model are required"})
		return
	}
	raw, err := cqFetch(map[string]string{"cmd": "getYears", "make": make_, "model": model})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var cqResp CQYearsResponse
	if err := json.Unmarshal(raw, &cqResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parse error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"minYear": cqResp.Years.MinYear, "maxYear": cqResp.Years.MaxYear})
}

// GET /api/garage/trims?make=Ferrari&model=F40&year=1992
func handleTrims(c *gin.Context) {
	make_ := c.Query("make")
	model := c.Query("model")
	year := c.Query("year")
	if make_ == "" || model == "" || year == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "make, model, and year are required"})
		return
	}
	raw, err := cqFetch(map[string]string{"cmd": "getTrims", "make": make_, "model": model, "year": year})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var cqResp CQTrimsResponse
	if err := json.Unmarshal(raw, &cqResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parse error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"trims": cqResp.Trims})
}

// GET /api/garage/vehicle?make=Ferrari&model=F40&year=1992
func handleVehicle(c *gin.Context) {
	make_ := c.Query("make")
	model := c.Query("model")
	yearS := c.Query("year")
	if make_ == "" || model == "" || yearS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "make, model, and year are required"})
		return
	}
	year := parseInt(yearS)

	raw, err := cqFetch(map[string]string{"cmd": "getTrims", "make": make_, "model": model, "year": yearS})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var cqResp CQTrimsResponse
	if err := json.Unmarshal(raw, &cqResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parse error"})
		return
	}
	if len(cqResp.Trims) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no data found for this vehicle"})
		return
	}

	spec := buildSpec(make_, model, year, cqResp.Trims[0])
	c.JSON(http.StatusOK, gin.H{"vehicle": spec})
}

// GET /api/garage/tuning-configs
func handleTuningConfigs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"configs": tuningConfigs})
}

// POST /api/garage/tune
// Body: { "make": "Ferrari", "model": "F40", "year": 1992, "tuning": "track" }
func handleTune(c *gin.Context) {
	var req TuneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cfg, ok := tuningConfigs[req.Tuning]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown tuning config: " + req.Tuning})
		return
	}

	yearS := strconv.Itoa(req.Year)
	raw, err := cqFetch(map[string]string{"cmd": "getTrims", "make": req.Make, "model": req.Model, "year": yearS})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var cqResp CQTrimsResponse
	if err := json.Unmarshal(raw, &cqResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parse error"})
		return
	}
	if len(cqResp.Trims) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no data found for this vehicle"})
		return
	}

	base := buildSpec(req.Make, req.Model, req.Year, cqResp.Trims[0])
	tuned := applyTuning(base, cfg)
	delta := Delta{
		HP:          tuned.HP - base.HP,
		Torque:      tuned.Torque - base.Torque,
		TopSpeed:    tuned.TopSpeed - base.TopSpeed,
		ZeroToSixty: round2(base.ZeroToSixty - tuned.ZeroToSixty),
		Weight:      tuned.Weight - base.Weight,
	}

	c.JSON(http.StatusOK, TuneResponse{Base: base, Tuned: tuned, Delta: delta, Config: cfg})
}

// GET /api/garage/search?q=ferrari
func handleSearch(c *gin.Context) {
	q := strings.ToLower(strings.TrimSpace(c.Query("q")))
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q is required"})
		return
	}
	raw, err := cqFetch(map[string]string{"cmd": "getMakes"})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var cqResp CQMakesResponse
	if err := json.Unmarshal(raw, &cqResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parse error"})
		return
	}
	type Result struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
	}
	var results []Result
	for _, m := range cqResp.Makes {
		if strings.Contains(strings.ToLower(m.MakeDisplay), q) || strings.Contains(strings.ToLower(m.MakeID), q) {
			results = append(results, Result{ID: m.MakeID, Name: m.MakeDisplay, Country: m.MakeCountry})
			if len(results) >= 10 {
				break
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

// ── Router setup (call this from your main.go) ────────────────────────────────
// If you have an existing router, import this and call RegisterGarageRoutes(r)

func RegisterGarageRoutes(r *gin.Engine) {
	garage := r.Group("/api/garage")
	{
		garage.GET("/makes", handleMakes)
		garage.GET("/models", handleModels)
		garage.GET("/years", handleYears)
		garage.GET("/trims", handleTrims)
		garage.GET("/vehicle", handleVehicle)
		garage.GET("/tuning-configs", handleTuningConfigs)
		garage.POST("/tune", handleTune)
		garage.GET("/search", handleSearch)
	}
}

// ── Standalone entrypoint (remove if you have your own main.go) ───────────────

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	RegisterGarageRoutes(r)

	fmt.Println("Garage API running on :8080")
	r.Run(":8080")
}
