import { describe, expect, test } from "@jest/globals"
import { FeedService } from "./FeedService"

describe("FeedService", () => {
  test("addEvent modifies cache", async () => {
    const service = new FeedService()
    const feedId = "gvmt.rcks"
    const feed = await service.getFeedById(feedId)
    await service.addEvent(feedId, { id: 1, at: "1980", title: "Test" })

    expect(feed.events.pop()).toEqual({ id: 1, at: "1980", title: "Test" })
  })
})
