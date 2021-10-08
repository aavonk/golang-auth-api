import { useState } from 'react';
import Box from '../../common/Box';
import { Button } from '../../common/Button';
import Card from '../../common/Card';
import Input from '../../common/Input';
import styles from './register.module.css';

function RegisterPage() {
  const [disabled, setDisabled] = useState(true);
  const [values, setValues] = useState({
    email: '',
    firstName: '',
    lastName: '',
    password: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      [e.target.name]: e.target.value,
    });
  };
  return (
    <div className="flex flex-col w-full h-screen items-center">
      <div className="flex flex-row lg:justify-between justify-center lg:mx-5 lg:my-0 pt-14  mx-auto lg:w-1080 ">
        {/* styles.wrapper ^ */}
        <div className="flex flex-col w-96 md:w-auto lg:flex-auto ">
          {/* styles.left ^ */}
          <Card>
            <Box padding={{ bottom: 12 }}>
              <h2 className="text-2xl font-medium text-primary-dark ls-card-title">
                Create your ~App name~ account
              </h2>
            </Box>
            <Box margin={{ top: 20, bottom: 32 }}>
              <Box margin={{ bottom: 12 }}>
                <label htmlFor="firstNameInput">First name</label>
              </Box>
              <Input
                value={values.firstName}
                onChange={handleChange}
                name="firstName"
                autoComplete="given-name"
                id="firstNameInput"
              />
            </Box>
            <Box margin={{ top: 20, bottom: 32 }}>
              <Box margin={{ bottom: 12 }}>
                <label htmlFor="lastNameInput">Last name</label>
              </Box>
              <Input
                value={values.lastName}
                onChange={handleChange}
                name="lastName"
                autoComplete="family-name"
                id="lastNameInput"
              />
            </Box>
            <Box margin={{ top: 20, bottom: 32 }}>
              <Box margin={{ bottom: 12 }}>
                <label htmlFor="emailInput">Email</label>
              </Box>
              <Input
                value={values.email}
                onChange={handleChange}
                name="email"
                autoComplete="email"
                id="emailInput"
              />
            </Box>
            <Box margin={{ top: 20, bottom: 32 }}>
              <Box margin={{ bottom: 12 }}>
                <label htmlFor="passwordInput">Password</label>
              </Box>
              <Input
                value={values.password}
                onChange={handleChange}
                name="password"
                id="passwordInput"
                type="password"
                autoComplete="new-password"
              />
            </Box>
            <Button type="submit" disabled={disabled} fullWidth>
              Create account
            </Button>
          </Card>
        </div>
        <div className="hidden lg:flex lg:flex-col lg:flex-auto">
          {/* Styles.right ^ */}
          <h3>RightSide!!!</h3>
        </div>
      </div>
    </div>
  );
}

export default RegisterPage;
