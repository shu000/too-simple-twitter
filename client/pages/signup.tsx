import Head from 'next/head';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import styles from '../styles/Home.module.css';
import { myfetch } from '../util/myfetch';

export default function Home() {
  const [message, setMessage] = useState('');

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm();

  // ユーザ登録
  const onSubmit = async (data: any) => {
    try {
      await myfetch('POST', '/api/v1/user', {
        user_id: data.user_id,
        name: data.name,
        password: data.password,
      });
    } catch (e) {
      setMessage(e.message);
    }

    setMessage('登録しました');
  };

  return (
    <div className={styles.container}>
      <Head>
        <title>Clone Of Twitter</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <h3>ユーザ登録</h3>
        <form onSubmit={handleSubmit(onSubmit)}>
          <input {...register('user_id', { required: true, maxLength: 80 })} placeholder="ユーザID" />
          {errors.user_id && <span>80字以内です</span>}

          <input {...register('name', { required: true, maxLength: 80 })} placeholder="ユーザ名" />
          {errors.user_id && <span>80字以内です</span>}

          <input {...register('password', { required: true })} placeholder="パスワード" />
          {errors.user_id && <span>入力必須です</span>}

          <input type="submit" value="登録" />
        </form>
        <p>{message}</p>
      </main>
    </div>
  );
}
