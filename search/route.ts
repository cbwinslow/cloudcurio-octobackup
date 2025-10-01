import { semanticSearch } from '@/lib/search'
export const runtime = 'edge'
export async function GET(req: Request) {
  const { searchParams } = new URL(req.url)
  const q = searchParams.get('q') || ''
  if (!q) return Response.json({ results: [] })
  const results = await semanticSearch(q)
  return Response.json({ results })
}
