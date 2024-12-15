// @ts-check

/**
 * @param {import('aws-lambda').APIGatewayProxyEventV2} event
 * @param {import('aws-lambda').Context} context
 * @returns {Promise<import('aws-lambda').APIGatewayProxyResultV2>}
 */
export const handler = async (event, context) => {
  console.log("Lambda received event:", event);

  return {
    statusCode: 200,
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      message: "Hello from Lambda!",
      event,
      context,
    }),
  };
};
