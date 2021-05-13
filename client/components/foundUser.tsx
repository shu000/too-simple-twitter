import { FC } from 'react';
import { User } from '../model/user';
import styles from '../styles/SearchResult.module.css';

interface Props {
  user: User;
  isFollowing: boolean;
  onFollowClick: (user: User) => void;
}

export const FoundUser: FC<Props> = ({ user, isFollowing, onFollowClick }) => {
  return (
    <div className={styles.foundUser}>
      <p className={styles.foundUser_userName}>
        {user.name}
        <span className={styles.foundUser_userId}>@{user.user_id}</span>
      </p>
      {isFollowing && <p>フォロー済み</p>}
      {!isFollowing && (
        <button type="button" onClick={() => onFollowClick(user)}>
          フォローする
        </button>
      )}
    </div>
  );
};
