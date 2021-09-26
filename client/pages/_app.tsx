import "../styles/globals.css";
import type { AppProps } from "next/app";
import { AuthProvider } from "context/auth";
import { NextPage } from "next";
import Head from "next/head";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      {/* {Component.requiresAuth && (
        <Head>
          <script
            // If no token is found, redirect inmediately
            dangerouslySetInnerHTML={{
              __html: `if(!document.cookie || document.cookie.indexOf('auth-session') === -1)
            {location.replace(
              "/login?next=" +
                encodeURIComponent(location.pathname + location.search)
            )}
            else {document.documentElement.classList.add("render")}`,
            }}
          />
        </Head>
      )} */}
      {/* <AuthProvider>
        <Component {...pageProps} />
      </AuthProvider> */}

      <Component {...pageProps} />
    </>
  );
}
export default MyApp;
