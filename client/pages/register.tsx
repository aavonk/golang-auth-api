import { NextPage } from "next";
import Card from "@common/Card";
import styles from "../styles/register.module.css";
import Box from "@common/Box";
import React, { useState } from "react";
import Input from "@common/Input";
import { Button } from "@common/Button";

const RegisterPage: NextPage = () => {
  const [disabled, setDisabled] = useState(true);
  const [values, setValues] = useState({
    email: "",
    firstName: "",
    lastName: "",
    password: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      [e.target.name]: e.target.value,
    });
  };
  return (
    <div className={styles.root}>
      <div className={styles.wrapper}>
        <div className={styles.left}>
          <Card>
            <Box padding={{ bottom: 12 }}>
              <h2 className="title-card">Create your ~App name~ account</h2>
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
        <div className={styles.right}>
          <h3>RightSide!!!</h3>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;
