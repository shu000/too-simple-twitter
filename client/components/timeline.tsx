import { FC } from 'react';
import { Tweet } from '../model/tweet';
import styles from '../styles/Timeline.module.css';
import { SingleTweet } from './singleTweet';

interface Props {
  tweets: Tweet[];
}

export const Timeline: FC<Props> = ({ tweets }) => {
  if (tweets.length === 0) {
    return (
      <div>
        <h2>タイムライン</h2>
        <p>ツイートがありません。</p>
        <p>ツイートするか、誰かをフォローしましょう。</p>
      </div>
    );
  }

  return (
    <div>
      <h2>タイムライン</h2>
      <div className={styles.tweets}>
        {tweets.map((tweet) => (
          <SingleTweet key={tweet.tweet_id} tweet={tweet} />
        ))}
      </div>
    </div>
  );
};
