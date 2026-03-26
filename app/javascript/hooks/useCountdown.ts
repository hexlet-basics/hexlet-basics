import { useEffect, useState } from "react";

type CountdownState = {
  isRunning: boolean;
  remainingSeconds: number;
};

const millisecondsPerSecond = 1000;

function getRemainingSeconds(expiryTime: number) {
  const remainingMilliseconds = Math.max(expiryTime - Date.now(), 0);

  return Math.ceil(remainingMilliseconds / millisecondsPerSecond);
}

export default function useCountdown(expiryTimestamp: Date): CountdownState {
  const expiryTime = expiryTimestamp.getTime();
  const [remainingSeconds, setRemainingSeconds] = useState(() =>
    getRemainingSeconds(expiryTime),
  );

  useEffect(() => {
    setRemainingSeconds(getRemainingSeconds(expiryTime));

    const intervalId = setInterval(() => {
      setRemainingSeconds(getRemainingSeconds(expiryTime));
    }, millisecondsPerSecond);

    return () => {
      clearInterval(intervalId);
    };
  }, [expiryTime]);

  return {
    isRunning: remainingSeconds > 0,
    remainingSeconds,
  };
}
