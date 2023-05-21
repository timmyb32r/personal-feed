import { makeAutoObservable } from "mobx"
import { sleep } from "../lib/delays"
import { Feed, FeedHead } from "../models/Feed"
import { FeedEvent } from "../models/FeedEvent"


const MOCK_FEEDS: Feed[] = [
  {
    id: "gvmt.rcks",
    title: "The Truth",
    events: [
      {
        id: 1,
        at: "1986-12-21",
        title: "Florida woman found drunk in refrigirator",
        description: "Authorities are confused... but not surprised.",
      },
      {
        id: 2,
        at: "2017-05-06",
        title: "Refrigirator prices crush real estate market",
        description: "Boomers are to blame again.",
      },
    ],
  },
  {
    id: "itsfreezinghere",
    title: "North Pole Digest",
    events: [
      {
        id: 3,
        at: "1200-09-10",
        title: "Nothing happened",
      },
      {
        id: 4,
        at: "1201-03-01",
        title: "One penguin fell in a funny way",
        description: "But nobody noticed.",
      },
      {
        id: 5,
        at: "1761-12-10",
        title: "Nothing happened in a while",
      },
    ],
  },
]

export class FeedService {
  feedsCache: { [id: string]: Feed } = {}

  constructor() {
    makeAutoObservable(this)
  }

  async listFeeds(): Promise<FeedHead[]> {
    //---------------------------------------------------------
    await sleep(200)
    const response = await fetch("/api/source_ids", {
      method: 'GET',
    });
    if (!response.ok) {
      throw new Error();
    }
    if (response.body !== null) {
      let result = response.json()
      console.log(Date.now(), "timmyb32rQQQ:listFeeds:", result)
      return result
    }
    //---------------------------------------------------------
  }

  async getFeedById(id: string): Promise<Feed> {
    if (id == undefined) {
      return null
    }
    if (!this.feedsCache[id]) {
      await sleep(200)
      //---------------------------------------------------------
      const response = await fetch(`/api/source_id/${id}`, {
        method: 'GET',
      });
      if (!response.ok) {
        throw new Error();
      }
      if (response.body !== null) {
        let result = response.json()
        console.log(Date.now(), "timmyb32rQQQ:getFeedById:", result)
        return result
      }
      //---------------------------------------------------------
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
