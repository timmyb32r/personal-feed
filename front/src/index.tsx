import React from "react"
import { createRoot } from "react-dom/client"
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom"
import { FeedsPage } from "./containers/FeedsPage"
import { HomePage } from "./containers/HomePage"
import "./index.css"

const router = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />,
  },
  {
    path: "/feeds",
    element: <FeedsPage />,
  },
])

const root = createRoot(document.getElementById("app"))

root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
)
