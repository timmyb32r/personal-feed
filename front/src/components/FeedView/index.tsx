import React from "react"
import { Feed } from "../../models/Feed"
import style from "./style.module.css"

export function FeedView(props: { feed: Feed }) {
  return <>
    <h2>{props.feed.title}</h2>
    <hr />
    <ol className={style.events}>
      {props.feed.events.map((event) => {
        return <li className={style.event} key={event.id}>
          <h3>{event.title}</h3>
          <div>{new Date(event.at).toLocaleDateString()}</div>
          {event.description ? <p>{event.description}</p> : null}
        </li>
      })}
    </ol>
  </>
}
