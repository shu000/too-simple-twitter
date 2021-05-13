export interface User {
  user_id: string;
  name: string;
}

export interface UserProfile extends User {
  following: User[];
  followed: User[];
}
