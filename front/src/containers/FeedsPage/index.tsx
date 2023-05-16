import { observer } from "mobx-react-lite"
import React, { useCallback, useEffect, useState } from "react"
import { useSearchParams } from "react-router-dom"
import { FeedView } from "../../components/FeedView"
import { FeedsList } from "../../components/FeedsList"
import { Feed, FeedHead } from "../../models/Feed"
import { feedService } from "../../services/FeedService"
import style from "./style.module.css"

export const FeedsPage = observer(() => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [activeFeedId, setActiveFeedId] = useState<string>(searchParams["selected"])
  const [activeFeed, setActiveFeed] = useState<Feed>()
  const [feedHeads, setFeedHeads] = useState<FeedHead[]>()

  useEffect(() => {
    feedService.listFeeds().then((heads) => {
      setFeedHeads(heads)
    })
  })

  useEffect(() => {
    feedService.getFeedById(activeFeedId).then((feed) => {
      setActiveFeed(feed)
    })
  }, [activeFeedId])

  const handleSelect = useCallback((feedId: string) => {
    setSearchParams({
      ...searchParams,
      selected: feedId,
    })
    setActiveFeedId(feedId)
  }, [setSearchParams, setActiveFeedId])

  let counter = 0
  const handleAddEvent = useCallback(() => {
    feedService.addEvent(activeFeedId, {
      id: counter,
      at: new Date().toISOString(),
      title: `Someone clicked The Button${counter++ > 0 ? " again" : ""}`,
      description: "But nobody cares...",
    })
  }, [activeFeedId])

  const feedBlock = () => <>
    <FeedView feed={activeFeed}/>
    <button onClick={handleAddEvent}>The Button</button>
  </>

return <div className={style.feedsPage}>
    <div className={style.feedsList}>
      {feedHeads ? <FeedsList activeId={activeFeedId} feeds={feedHeads} onSelect={handleSelect}/> : <div>Loading feeds...</div>}
    </div>
    <div className={style.feed}>
      {activeFeed ? feedBlock() : <div>No active feed...</div>}
    </div>
  </div>
})
