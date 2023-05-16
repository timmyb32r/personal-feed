import React from "react"
import { Link } from "react-router-dom"

export function HomePage() {
  return <ul>
    <li><Link to="/">Home</Link></li>
    <li><Link to="/feeds">Feeds</Link></li>
  </ul>
}
