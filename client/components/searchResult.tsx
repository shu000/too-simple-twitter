import { FC } from 'react';
import { User } from '../model/user';
import styles from '../styles/SearchResult.module.css';
import { FoundUser } from './foundUser';

interface Props {
  users: User[];
  followings: User[];
  onFollowClick: (user: User) => void;
}

export const SerachResult: FC<Props> = ({ users, followings, onFollowClick }) => {
  if (users.length === 0) {
    return (
      <div className={styles.SerachResult}>
        <p>検索結果が０件です。</p>
      </div>
    );
  }

  return (
    <div className={styles.searchResult}>
      {users.map((user) => {
        const isFollowing = followings.some((following) => following.user_id === user.user_id);
        return <FoundUser key={user.user_id} user={user} isFollowing={isFollowing} onFollowClick={onFollowClick} />;
      })}
    </div>
  );
};
