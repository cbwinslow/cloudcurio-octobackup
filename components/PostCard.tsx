import { Post } from 'contentlayer/generated'
export default function PostCard({ post }:{ post: Post }){
  return (
    <a href={`/blog/${post.slug}`} className="block rounded bg-cloudcurio-surface p-4 hover:shadow-glow">
      <div className="font-semibold text-cloudcurio-mint">{post.title}</div>
      <p className="text-xs opacity-80 mt-1">{post.summary}</p>
      <div className="text-xs mt-2 opacity-70">{new Date(post.date).toLocaleDateString()}</div>
    </a>
  )
}
