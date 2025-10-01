export default function Footer(){
  return (
    <footer className="mt-16">
      <div className="max-w-6xl mx-auto p-6 text-xs opacity-70">
        © {new Date().getFullYear()} Cloudcurio · Built with Next.js on Cloudflare
      </div>
    </footer>
  )
}
