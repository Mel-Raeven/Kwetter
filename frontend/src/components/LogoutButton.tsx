import React from 'react';

const LogoutButton: React.FC = () => {
  // Clear local storage
  sessionStorage.removeItem('access_token');
  sessionStorage.removeItem('id_token');

  // Clear HTTP-only cookies
  document.cookie = 'accessToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
  document.cookie = 'idToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
  document.cookie = 'refreshToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';

  // Redirect to Cognito sign-out URL
  const url = 'https://kwetter.auth.eu-central-1.amazoncognito.com/logout?client_id=5vs4anjqtme5f93037vv30fmp0&logout_uri=http://localhost:5173/';

  return (
    <button onClick={() => window.location.href = url}>Logout</button>
  );
}

export default LogoutButton;
