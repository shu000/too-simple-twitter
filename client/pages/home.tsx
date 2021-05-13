import Head from 'next/head';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { SerachResult } from '../components/searchResult';
import { Timeline } from '../components/timeline';
import { Tweet } from '../model/tweet';
import { User, UserProfile } from '../model/user';
import styles from '../styles/Home.module.css';
import { myfetch } from '../util/myfetch';

export default function Home() {
  const router = useRouter();

  const tweetForm = useForm();
  const searchForm = useForm();

  const [message, setMessage] = useState('');
  const [me, setMe] = useState<UserProfile | null>(null);
  const [tweets, setTweets] = useState<Tweet[]>([]);
  const [users, setUsers] = useState<User[]>([]);

  // 画面読み込み時にログイン中ユーザのプロフィールを取得
  useEffect(() => {
    (async () => {
      try {
        const result = await myfetch('GET', '/api/v1/user/me');
        setMe(result);
      } catch (e) {
        router.push('/app');
      }
    })();
  }, [setMe]);

  // 画面読み込み時にタイムラインを取得
  useEffect(() => {
    (async () => {
      try {
        const result = await myfetch('GET', '/api/v1/tweet');
        setTweets(result);
      } catch (e) {
        setMessage(e.message);
      }
    })();
  }, [setTweets, me?.following]); // フォロー時に反映させる雑な手段

  // ツイートクリックハンドラ
  const onTweetSubmit = async (data: any) => {
    try {
      const result = await myfetch('POST', '/api/v1/tweet', {
        contents: data.contents,
      });
      setTweets([result, ...tweets]);
      tweetForm.reset();
    } catch (e) {
      setMessage(e.message);
    }
  };

  // ユーザ検索クリックハンドラ
  const onSearchSubmit = async (data: any) => {
    try {
      const result = await myfetch('GET', `/api/v1/user?user_id=${data.user_id}`);
      setUsers(result);
    } catch (e) {
      setMessage(e.message);
    }
  };

  // フォローするクリックハンドラ
  const onFollowClick = async (user: User) => {
    try {
      await myfetch('PATCH', '/api/v1/user/follow', { target_user_id: user.user_id });
      if (me) {
        setMe({
          ...me,
          following: [...me.following, user],
        });
      }
    } catch (e) {
      setMessage(e.message);
    }
  };

  return (
    <div className={styles.container}>
      <Head>
        <title>Clone Of Twitter</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <p>{message}</p>
      <main className={styles.tweColumn}>
        <div className={styles.tweColumn_innerWrapper}>
          <Timeline tweets={tweets} />
        </div>

        <div className={styles.tweColumn_innerWrapper}>
          <h3>ツイートする</h3>
          <form onSubmit={tweetForm.handleSubmit(onTweetSubmit)}>
            <textarea {...tweetForm.register('contents', { required: true, maxLength: 140 })} placeholder="内容" />
            {tweetForm.formState.errors.contents && <span>ツイートは1文字以上140文字以下です</span>}
            <input type="submit" value="ツイート" />
          </form>

          <h3>ユーザ検索する</h3>
          <form onSubmit={searchForm.handleSubmit(onSearchSubmit)}>
            <input {...searchForm.register('user_id', { required: true })} placeholder="@ユーザID" />
            <input type="submit" value="検索" />
            <SerachResult users={users} followings={me ? me.following : []} onFollowClick={onFollowClick} />
          </form>
        </div>
      </main>
    </div>
  );
}
