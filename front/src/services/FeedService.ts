import { makeAutoObservable } from "mobx"
import { sleep } from "../lib/delays"
import { Feed, FeedHead } from "../models/Feed"
import { FeedEvent } from "../models/FeedEvent"
import { wrapURI } from "./core"


export class FeedService {
  feedsCache: { [id: string]: Feed } = {}
  headsCache?: FeedHead[] = undefined

  constructor() {
    makeAutoObservable(this)
  }

  async listFeeds(): Promise<FeedHead[]> {
    if (!this.headsCache) {
      const response = await fetch(wrapURI("/api/source_ids"), {
        method: 'GET',
      });
      if (!response.ok) {
        throw new Error();
      }
      this.headsCache = await response.json()
    }
    return this.headsCache
  }

  async getFeedById(id: string): Promise<Feed> {
    if (id == undefined) {
      return null
    }
    if (!this.feedsCache[id]) {
      const response = await fetch(wrapURI(`/api/source_id/${id}`), {
        method: 'GET',
        credentials: "omit",
      });
      if (!response.ok) {
        throw new Error();
      }
      this.feedsCache[id] = await response.json()
    }

    return this.feedsCache[id]
  }

  async addEvent(feedId: string, event: FeedEvent): Promise<void> {
    await sleep(200)

    const feed = await this.getFeedById(feedId)
    feed.events.push(event)
  }
}

export const feedService = new FeedService()
