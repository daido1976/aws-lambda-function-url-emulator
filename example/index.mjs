export const handler = async (event) => {
  console.log("Lambda received event:", event);

  const response = {
    statusCode: 200,
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      message: "Hello from Lambda!",
      input: event,
    }),
  };

  return response;
};
