import { NextPage } from "next";
import styles from "../styles/login.module.css";
import Card from "@common/Card";
import Input from "@common/Input";
import React, { useState } from "react";
import Box from "@common/Box";
import { Button } from "@common/Button";

const LoginPage: NextPage = () => {
  const [values, setValues] = useState({
    email: "",
    password: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log(values);
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
              <Button fullWidth type="submit">
                Continue
              </Button>
            </form>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
