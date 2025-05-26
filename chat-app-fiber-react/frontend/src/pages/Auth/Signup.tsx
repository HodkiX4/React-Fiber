import { FormEvent } from 'react';
import styles from './auth.module.css';
import { useAuth } from '../../hooks/useAuth';
import { ISignupRequest } from '../../models/auth.model';
import { Link } from 'react-router-dom';

const Signup = () => {
  const { signup } = useAuth();
  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const values = Object.fromEntries(formData) as unknown as ISignupRequest
    await signup(values)
  }

  return (
    <div className={styles.container}>
      <form className={styles.signupForm} onSubmit={handleSubmit}>
        <div className={styles.header}>
          <h2>Signup</h2>
          <span>Create a new account</span>
        </div>
        <div className={styles.inputGroup}>
          <input type="text" id="username" placeholder="" name='username' />
          <label htmlFor="username">Username</label>
        </div>
        <div className={styles.inputGroup}>
          <input type="text" id="email" placeholder="" name='email' />
          <label htmlFor="email">Email</label>
        </div>
        <div className={styles.inputGroup}>
          <input type="password" id="password" placeholder="" name='password' />
          <label htmlFor="password">Password</label>
        </div>
        <div className={styles.inputGroup}>
          <input type="password" id="confirm-password" placeholder="" name='passwordAgain' />
          <label htmlFor="confirm-password">Confirm Password</label>
        </div>
        <button className={styles.submit} type="submit">
          Submit
        </button>
        <div className={styles.footer}>
          <p>Already have an account? Log in <Link to="/login">here</Link></p>
        </div>
      </form>
    </div>
  );
}

export default Signup;