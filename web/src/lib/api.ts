// API client for the opensource-curator Go backend

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface Library {
  id: string;
  name: string;
  registry: string;
  packageName: string;
  githubRepo: string;
  description: string;
  latestVersion: string;
  latestVersionDate?: string;
  deprecated: boolean;
  score?: {
    overall: number;
    breakdown: ScoreBreakdown;
    version: string;
  };
  alternatives?: Alternative[];
}

export interface ScoreBreakdown {
  maintenanceHealth: number;
  apiClarity: number;
  docQuality: number;
  securityPosture: number;
  communitySignal: number;
  deprecationSafety: number;
}

export interface Alternative {
  name: string;
  registry: string;
  packageName: string;
  overallScore: number;
  relationship: string;
  reason: string;
}

export interface Category {
  id: string;
  slug: string;
  name: string;
  description: string;
  libraryCount?: number;
}

export interface Action {
  rel: string;
  href: string;
  description?: string;
}

export interface Envelope<T> {
  data: T;
  error?: string;
  next_actions?: Action[];
}

export interface Recommendation {
  rank: number;
  library: Library;
  matchReason: string;
  preferenceMatch?: string;
}

export interface RecommendResponse {
  data: Recommendation[];
  query: {
    task: string;
    lang?: string;
    prefer?: string;
    matchedCategories: string[];
  };
  next_actions?: Action[];
}

async function fetchAPI<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    next: { revalidate: 3600 },
  });
  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`);
  }
  return res.json();
}

export async function getCategories(): Promise<Envelope<Category[]>> {
  return fetchAPI("/v1/categories");
}

export async function getCategory(slug: string): Promise<Envelope<Category & { libraries: Library[] }>> {
  return fetchAPI(`/v1/categories/${slug}`);
}

export async function getLibrary(id: string): Promise<Envelope<Library>> {
  return fetchAPI(`/v1/libraries/${id}`);
}

export async function getLibraryBySlug(registry: string, packageName: string): Promise<Envelope<Library>> {
  return fetchAPI(`/v1/libraries/${encodeURIComponent(registry)}/${encodeURIComponent(packageName)}`);
}

export async function searchLibraries(query: string): Promise<Envelope<Library[]>> {
  return fetchAPI(`/v1/search?q=${encodeURIComponent(query)}`);
}

export async function recommend(task: string, prefer?: string): Promise<RecommendResponse> {
  let path = `/v1/recommend?task=${encodeURIComponent(task)}`;
  if (prefer) path += `&prefer=${encodeURIComponent(prefer)}`;
  return fetchAPI(path);
}

export async function getHealth(): Promise<Envelope<{ status: string; timestamp: string }>> {
  return fetchAPI("/v1/health");
}

export async function getScoringWeights(): Promise<Envelope<ScoreBreakdown & { version: string }>> {
  return fetchAPI("/v1/scoring/weights");
}
