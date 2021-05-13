import Head from 'next/head';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import styles from '../styles/Home.module.css';
import { myfetch } from '../util/myfetch';

export default function Login() {
  const [message, setMessage] = useState('');

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm();

  const onSubmit = async (data: any) => {
    try {
      const result = await myfetch('POST', '/api/v1/auth/login', {
        user_id: data.user_id,
        password: data.password,
      });

      localStorage.setItem('MyToken', result.token);
      window.location.href = '/app/home';
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

      <main className={styles.main}>
        <form onSubmit={handleSubmit(onSubmit)}>
          <input {...register('user_id', { required: true })} placeholder="ユーザID" />
          {errors.user_id && <span>入力必須です</span>}

          <input {...register('password', { required: true })} placeholder="パスワード" />
          {errors.password && <span>入力必須です</span>}

          <input type="submit" value="ログイン" />
        </form>
        <p>{message}</p>
      </main>
    </div>
  );
}
