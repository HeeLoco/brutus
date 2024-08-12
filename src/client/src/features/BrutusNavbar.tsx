import Logo from '/logo.svg'

export function BrutusNavbar() {
  return (
    <nav className="flex flex-row gap-3 justify-start items-center shadow-lg p-2">
      <img src={Logo} className="h-10 w-10"></img>
      Brutus
    </nav>
  )
}

export default BrutusNavbar;