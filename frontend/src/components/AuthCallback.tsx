import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const AuthCallback: React.FC = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');

    if (code) {
      fetchTokens(code).then(() => {
        console.log('Tokens fetched and stored successfully');
        navigate('/'); // Redirect to home page after successful authentication
      }).catch((error) => {
        console.error('Error fetching tokens:', error);
      });
    } else {
      console.error('Authorization code not found');
    }
  }, [navigate]);

  const fetchTokens = async (code: string) => {
    try {
      const clientId = '5vs4anjqtme5f93037vv30fmp0';
      const redirectUri = 'http://localhost:5173/auth/callback';
      const tokenUrl = 'https://kwetter.auth.eu-central-1.amazoncognito.com/oauth2/token';
      const body = `grant_type=authorization_code&client_id=${clientId}&code=${code}&redirect_uri=${redirectUri}`;

      const response = await fetch(tokenUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body,
      });

      if (!response.ok) {
        throw new Error('Failed to fetch tokens');
      }

      const data = await response.json();
      console.log('Fetched token data:', data);

      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('id_token', data.id_token);
    } catch (error) {
      console.error('Error during fetchTokens:', error);
    }
  };

  return <div>Loading...</div>;
};

export default AuthCallback;
