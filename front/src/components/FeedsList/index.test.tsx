import { describe, expect, test } from "@jest/globals"
import React from "react"
import renderer from "react-test-renderer"
import { FeedsList } from "."
import { FeedHead } from "../../models/Feed"

const mockFeeds: FeedHead[] = [
  { id: "1", title: "Test 1" },
  { id: "2", title: "Test 2" },
  { id: "3", title: "Test 3" },
]

describe("FeedsList", () => {
  test("Active item is highlighted", () => {
    const component = renderer.create(<FeedsList feeds={mockFeeds} activeId='2' />)

    expect(component.toJSON()).toMatchSnapshot()
  })
})
