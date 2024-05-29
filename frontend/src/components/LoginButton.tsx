import React from 'react';

const LoginButton: React.FC = () => {
  const clientId = '1vnld5vl7tocmecvbkjksmsf2r';
  const redirectUri = encodeURIComponent('http://localhost:5173/auth/callback');
  const cognitoDomain = 'https://kwetter.auth.eu-central-1.amazoncognito.com/oauth2/authorize';

  const loginUrl = `${cognitoDomain}?client_id=${clientId}&response_type=code&scope=email+openid+phone&redirect_uri=${redirectUri}`;

  return (
    <button onClick={() => window.location.href = loginUrl}>
      Login
    </button>
  );
};

export default LoginButton;
