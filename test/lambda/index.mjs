// @ts-check

/**
 * @param {import('aws-lambda').APIGatewayProxyEventV2} event
 * @param {import('aws-lambda').Context} context
 * @returns {Promise<import('aws-lambda').APIGatewayProxyResultV2>}
 */
export const handler = async (event, context) => {
  console.log("Lambda received event:", event);

  // Check if the request is for binary data
  const isBinary = event.queryStringParameters?.binary === "true";
  if (isBinary) {
    const binaryData = Buffer.from(
      "This is a test binary data",
      "utf-8"
    ).toString("base64");
    return {
      statusCode: 200,
      headers: {
        "Content-Type": "application/octet-stream",
      },
      body: binaryData,
      isBase64Encoded: true,
    };
  }

  // default response
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
