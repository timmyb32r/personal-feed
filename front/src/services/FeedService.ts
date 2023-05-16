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
    await sleep(200)

    return MOCK_FEEDS.map((x) => {
      return { title: x.title, id: x.id }
    })
  }

  async getFeedById(id: string): Promise<Feed> {
    if (!this.feedsCache[id]) {
      await sleep(200)

      this.feedsCache[id] = MOCK_FEEDS.find((feed) => {
        return feed.id === id
      })
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
