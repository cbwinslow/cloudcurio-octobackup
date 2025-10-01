'use client'
import { useEffect, useRef } from 'react'
type Props = { gistId: string; file?: string }
export default function GistEmbed({ gistId, file }: Props) {
  const ref = useRef<HTMLDivElement>(null)
  useEffect(() => {
    const s = document.createElement('script')
    s.src = `https://gist.github.com/${gistId}.js${file ? `?file=${file}` : ''}`
    s.async = true
    ref.current?.appendChild(s)
    return () => { if (ref.current) ref.current.innerHTML = '' }
  }, [gistId, file])
  return <div className="gist-embed" ref={ref} />
}
