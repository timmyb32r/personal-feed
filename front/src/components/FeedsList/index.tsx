import React from "react"
import { FeedHead } from "../../models/Feed"
import style from "./style.module.css"

export function FeedsList(props: {
  feeds: FeedHead[],
  activeId: string,
  onSelect?: (feedId: string) => void
}) {
  return <ol className={style.list}>
    {props.feeds.map((feed) => {
      const cls = props.activeId === feed.id ? style.active : ""
      return <li key={feed.id} className={cls} onClick={() => props.onSelect?.(feed.id)}>{feed.title}</li>
    })}
  </ol>
}
