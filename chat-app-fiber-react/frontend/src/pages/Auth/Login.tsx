import { FormEvent } from 'react';
import styles from './auth.module.css';
import { Link, useNavigate } from 'react-router-dom';
import { ILoginRequest } from '../../models/auth.model';
import { useAuth } from '../../hooks/useAuth';

function Login() {
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const values = Object.fromEntries(formData) as unknown as ILoginRequest;
    const success = await login(values)
    if(success) {
      navigate('/')
    }
  }

  return (
    <div className={styles.container}>
      <form className={styles.signupForm} onSubmit={handleSubmit}>
        <div className={styles.header}>
          <h2>Login</h2>
          <span>Log in to your account</span>
        </div>
        <div className={styles.inputGroup}>
          <input type="text" id="email" placeholder="" name='email' />
          <label htmlFor="email">Email</label>
        </div>
        <div className={styles.inputGroup}>
          <input type="password" id="password" placeholder="" name='password' />
          <label htmlFor="password">Password</label>
        </div>
        <button className={styles.submit} type="submit">
          Submit          
        </button>
        <div className={styles.footer}>
          <p>You don't have an account? Sign up <Link to="/signup">here</Link></p>
        </div>
      </form>
    </div>
  )
}

export default Login;
