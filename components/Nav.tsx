export default function Nav(){
  return (
    <header className="w-full">
      <nav className="max-w-6xl mx-auto p-4 flex items-center justify-between">
        <a href="/" className="font-bold text-cloudcurio-mint">Cloudcurio</a>
        <div className="flex gap-4 text-sm">
          <a href="/blog">Blog</a>
          <a href="/search">Search</a>
          <a href="/github">GitHub</a>
          <a href="/chat">Chat</a>
          <a href="/logs">Logs</a>
          <a href="/admin">Admin</a>
        </div>
      </nav>
    </header>
  )
}
