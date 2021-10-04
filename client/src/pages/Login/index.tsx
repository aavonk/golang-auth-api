import { useState } from "react";
import axios from "axios";
import Box from "../../common/Box";
import { Button } from "../../common/Button";
import Card from "../../common/Card";
import Input from "../../common/Input";
import styles from "./login.module.css";
import { Link } from "react-router-dom";

function LoginPage() {
  const [values, setValues] = useState({
    email: "",
    password: "",
  });

  const [error, setError] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const res = await axios.post(`/api/auth/signin`, values);
      console.log(res.data);
    } catch (err) {
      setError(true);
    }
  };
  return (
    <div className={styles.root}>
      <div className={styles.contentWrapper}>
        <div className={styles.cardWrapper}>
          <Card>
            <form className={styles.cardBody} onSubmit={handleSubmit}>
              <Box>
                <h2 className="title-card">Sign in to your account</h2>
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
                  size="large"
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
                  size="large"
                />
              </Box>
              <Box margin={{ top: -10, bottom: 20 }}>
                {error && (
                  <div className="text-color--red">
                    Incorrect email or password
                  </div>
                )}
              </Box>
              <Button fullWidth type="submit">
                Continue
              </Button>
            </form>
          </Card>
          <Box margin={{ top: 32, left: 20 }}>
            <p>
              Dont have an account?{" "}
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
