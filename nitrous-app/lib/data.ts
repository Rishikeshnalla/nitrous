import type { Event, Category, Journey, MerchItem, NavCard } from '@/types'

export const events: Event[] = [
  {
    id: '1',
    title: 'NASCAR Daytona 500',
    location: 'Daytona International Speedway · Florida',
    date: 'Feb 16, 2026',
    isLive: true,
    category: 'motorsport',
  },
  {
    id: '2',
    title: 'Dakar Rally — Stage 9',
    location: 'Al Ula → Ha'il · Saudi Arabia',
    date: 'Feb 10, 2026',
    time: '09:00 UTC',
    isLive: false,
    category: 'offroad',
  },
  {
    id: '3',
    title: 'World Dirt Track Championship',
    location: 'Knob Noster · Missouri, USA',
    date: 'Feb 21, 2026',
    isLive: true,
    category: 'motorsport',
  },
  {
    id: '4',
    title: 'Speed Boat Cup — Finals',
    location: 'Lake Como · Italy',
    date: 'Mar 2, 2026',
    time: '14:00 UTC',
    isLive: false,
    category: 'water',
  },
  {
    id: '5',
    title: 'Red Bull Skydive Series — Rd. 3',
    location: 'Interlaken Drop Zone · Switzerland',
    date: 'Mar 8, 2026',
    time: '11:30 UTC',
    isLive: false,
    category: 'air',
  },
  {
    id: '6',
    title: 'Crop Duster Air Racing',
    location: 'Bakersfield Airfield · California',
    date: 'Mar 14, 2026',
    time: '16:00 UTC',
    isLive: false,
    category: 'air',
  },
]

export const categories: Category[] = [
  {
    id: '1',
    name: 'MOTORSPORT',
    slug: 'motorsport',
    icon: '🏎️',
    liveCount: 24,
    description: 'NASCAR · F1 · Dirt · Rally',
    color: 'cyan',
  },
  {
    id: '2',
    name: 'WATER',
    slug: 'water',
    icon: '🌊',
    liveCount: 8,
    description: 'Speed Boats · Jet Ski · Surf',
    color: 'blue',
  },
  {
    id: '3',
    name: 'AIR & SKY',
    slug: 'air',
    icon: '🪂',
    liveCount: 5,
    description: 'Skydive · Air Race · Wing',
    color: 'purple',
  },
  {
    id: '4',
    name: 'OFF-ROAD',
    slug: 'offroad',
    icon: '🏔️',
    liveCount: 12,
    description: 'Dakar · Baja · Enduro',
    color: 'orange',
  },
]

export const journeys: Journey[] = [
  {
    id: '1',
    title: 'DAYTONA PIT CREW EXPERIENCE',
    category: 'MOTORSPORT · BEHIND THE SCENES',
    description: 'Go behind the wall at Daytona 500. Watch pit stops up close, meet the crew chiefs, and ride the pace car on track.',
    badge: 'EXCLUSIVE',
    slotsLeft: 12,
    date: 'Feb 16, 2026',
    price: 2400,
  },
  {
    id: '2',
    title: 'DAKAR DESERT CONVOY',
    category: 'RALLY · DESERT EXPEDITION',
    description: 'Ride a support vehicle through the Dakar stages. Sleep under the stars, eat with the team, and feel the dust.',
    badge: 'MEMBERS ONLY',
    slotsLeft: 6,
    date: 'Jan 18, 2027',
    price: 5800,
  },
  {
    id: '3',
    title: 'RED BULL TANDEM SKYDIVE',
    category: 'AIR · EXTREME SPORT',
    description: 'Jump with a Red Bull certified instructor at 15,000ft. Camera-equipped, full debrief, and a story you'll never forget.',
    badge: 'LIMITED',
    slotsLeft: 3,
    date: 'Mar 8, 2026',
    price: 1200,
  },
]

export const merchItems: MerchItem[] = [
  { id: '1', name: 'Team Hoodie', icon: '👕', price: 89, category: 'apparel' },
  { id: '2', name: 'NITROUS Cap', icon: '🧢', price: 42, category: 'apparel' },
  { id: '3', name: 'Racing Jacket', icon: '🏎️', price: 189, category: 'apparel' },
  { id: '4', name: 'Pit Watch', icon: '⌚', price: 249, category: 'accessories' },
  { id: '5', name: 'Gear Backpack', icon: '🎒', price: 120, category: 'accessories' },
  { id: '6', name: 'Drop Keychain', icon: '🏆', price: 28, category: 'collectibles' },
]

export const heroNavCards: NavCard[] = [
  { id: '1', label: 'ACCESS\nGARAGE', icon: '🚗', href: '/garage', color: 'grey', progress: 60 },
  { id: '2', label: 'ACCESS\nEVENT PASSES', icon: '🎫', href: '/passes', color: 'red', progress: 40 },
  { id: '3', label: 'ACCESS\nLIVE STREAMS', icon: '📺', href: '/live', color: 'cyan', progress: 75 },
  { id: '4', label: 'TEAMS', icon: '🏆', href: '/teams', color: 'orange', progress: 55 },
  { id: '5', label: 'JOURNEYS', icon: '🌍', href: '/journeys', color: 'blue', progress: 85 },
  { id: '6', label: 'MERCH', icon: '👕', href: '/merch', color: 'gold', progress: 45 },
]
