export const runtime = 'edge'

async function embed(text: string) {
  // @ts-ignore
  const res = await (globalThis as any).AI.run('@cf/baai/bge-base-en-v1.5', { text })
  return (res.data?.[0] ?? []) as number[]
}

export async function indexPost(doc: { slug:string; title:string; summary:string; body:string }) {
  // @ts-ignore
  const vec = (globalThis as any).VEC
  const embedding = await embed([doc.title, doc.summary, doc.body].join('\n'))
  await vec.upsert([{ id: doc.slug, values: embedding, metadata: { title: doc.title, summary: doc.summary } }])
}

export async function semanticSearch(q: string, k=8) {
  // @ts-ignore
  const vec = (globalThis as any).VEC
  const qv = await embed(q)
  const r = await vec.query({ topK: k, vector: qv, returnValues: false, includeMetadata: true })
  return r.matches as Array<{ id:string; score:number; metadata:{title:string; summary:string} }>
}
