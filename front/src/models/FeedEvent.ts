import { ISOTimestamp } from "./Timestamp"

export type FeedEvent = {
  id: number,
  at: ISOTimestamp
  title: string
  description?: string
}
