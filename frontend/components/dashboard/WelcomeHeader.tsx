// components/dashboard/WelcomeHeader.tsx
interface WelcomeHeaderProps {
  name: string;
}

const WelcomeHeader = ({ name }: WelcomeHeaderProps) => {
  return (
    <div>
      <h1 className="text-3xl font-bold text-dark-background">
        Welcome to the Access Control Manager, {name}! 👋
      </h1>
    </div>
  );
};

export default WelcomeHeader;