import { FeedEvent } from "./FeedEvent"

export type FeedHead = {
  id: string,
  title: string,
}

export type Feed = FeedHead & {
  events: FeedEvent[]
}
