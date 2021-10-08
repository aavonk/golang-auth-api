import { useState } from 'react';
import axios from 'axios';
import Box from '../../common/Box';
import { Button } from '../../common/Button';
import Card from '../../common/Card';
import Input from '../../common/Input';
import { Link } from 'react-router-dom';
import { useForm } from '../../hooks/useForm';

type LoginFields = {
  email: string;
  password: string;
};

function LoginPage() {
  const { values, errors, handleChange, handleSubmit } = useForm<LoginFields>({
    validations: {
      email: {
        required: {
          value: true,
          message: 'Please enter your email',
        },
      },
      password: {
        required: {
          value: true,
          message: 'Please enter your password',
        },
      },
    },
    initialValues: {
      email: '',
      password: '',
    },
    onSubmit: () => {
      login();
    },
  });
  const [networkError, setNetworkError] = useState(false);

  const login = async () => {
    try {
      const res = await axios.post(`/api/auth/signin`, values);
      console.log(res.status);
      //@ts-ignore
      if (res.data.error) {
        console.log('EROROROROROR');
      }
    } catch (err) {
      if (axios.isAxiosError(err)) {
        return setNetworkError(true);
      }

      alert('SOMETHING WENT WRONG');
    }
  };
  return (
    <div className="flex flex-col w-full h-screen items-center">
      <div className="flex flex-row w-full lg:w-1085 h-full justify-center items-center">
        <div className={'w-96 md:w-1/2'}>
          <Card>
            <form className="px-5" onSubmit={handleSubmit}>
              <Box>
                <h2 className=" text-2xl font-medium text-primary-dark ls-card-title ">
                  Sign in to your account
                </h2>
              </Box>
              <Box margin={{ top: 20, bottom: 32 }}>
                <Box margin={{ bottom: 12 }}>
                  <label htmlFor="emailInput" className="text-primary-reg">
                    Email
                  </label>
                </Box>
                <Input
                  value={values.email}
                  onChange={handleChange('email')}
                  name="email"
                  autoComplete="email"
                  id="emailInput"
                  size="large"
                />
                {errors.email && (
                  <p className="text-error" style={{ marginTop: 10 }}>
                    {errors.email}
                  </p>
                )}
              </Box>
              <Box margin={{ top: 20, bottom: 32 }}>
                <Box margin={{ bottom: 12 }}>
                  <label htmlFor="passwordInput">Password</label>
                </Box>
                <Input
                  value={values.password}
                  onChange={handleChange('password')}
                  name="password"
                  id="passwordInput"
                  type="password"
                  size="large"
                />
                {errors.password && (
                  <p className="text-error" style={{ marginTop: 10 }}>
                    {errors.password}
                  </p>
                )}
              </Box>
              <Box margin={{ top: -10, bottom: 20 }}>
                {networkError && (
                  <div className="text-error">Incorrect email or password</div>
                )}
              </Box>
              <Button fullWidth type="submit" className="mb-5 mt-10">
                Continue
              </Button>
            </form>
          </Card>
          <Box margin={{ top: 32, left: 20 }}>
            <p>
              Dont have an account?{' '}
              <Link to="/register" className="link">
                Sign up
              </Link>
            </p>
          </Box>
        </div>
      </div>
    </div>
  );
}

export default LoginPage;
