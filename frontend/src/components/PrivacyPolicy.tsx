// PrivacyPolicy.tsx
import { Container, Title, Text } from '@mantine/core';

const PrivacyPolicy: React.FC = () => {
  return (
    <Container>
      <Title order={1}>Privacy Policy</Title>
      <Text>
        Our application stores the messages created by the user. All user information, including the messages, is attached to their email address. We ensure that user data is kept confidential and secure. By using our service, you agree to the storage and use of your information as described in this policy.
      </Text>
    </Container>
  );
};

export default PrivacyPolicy;

