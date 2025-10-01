export const runtime = 'edge'
export async function POST(req: Request){
  const { text } = await req.json()
  // @ts-ignore
  const out = await (globalThis as any).AI.run('@cf/meta/llama-3.1-8b-instruct', { messages:[{role:'user', content:`Summarize for a blog sidebar: ${text}`}] })
  return Response.json({ summary: out.response })
}
