import { FC } from 'react';
import { Tweet } from '../model/tweet';
import styles from '../styles/Timeline.module.css';

interface Props {
  tweet: Tweet;
}

export const SingleTweet: FC<Props> = ({ tweet }) => {
  return (
    <div className={styles.singleTweet}>
      <p className={styles.singleTweet_userName}>
        {tweet.tweet_user_name}
        <span className={styles.singleTweet_userId}>@{tweet.tweet_user_id}</span>
      </p>
      <p>{tweet.contents}</p>
    </div>
  );
};
