import axios from "axios";

import { GetServerSideProps, GetServerSidePropsContext } from "next";

export function requireAuthentication(gssp: GetServerSideProps) {
  return async (ctx: GetServerSidePropsContext) => {
    const { req } = ctx;

    if (!req.headers.cookie) {
      return {
        redirect: {
          permanant: false,
          destination: "/login",
        },
      };
    }

    try {
      console.log("FETCHING DATA");
      const response = await axios.get("http://auth-api:7777/v1/user/me", {
        headers: req.headers,
      });
      console.log(response);
      console.log("FETCH RETURNED");
      if (!response.data.data) {
        console.log("REDIRECT FROM IF--");
        return {
          redirect: {
            permanant: false,
            destination: "/login",
          },
        };
      }
    } catch (error) {
      // Failure in the query or any error should fallback here
      // this route is possibly forbidden means the cookie is invalid
      // or the cookie is expired
      console.error(error);
      return {
        redirect: {
          permanant: false,
          destination: "/login",
        },
      };
    }

    return await gssp(ctx);
  };
}
