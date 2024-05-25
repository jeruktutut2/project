/** @type {import('next').NextConfig} */
const nextConfig = {
    // async rewrites() {
    //     return [
    //       {
    //         source: '/api',
    //         destination: 'http://localhost/api',
    //       },
    //     ]
    // },
    async rewrites() {
		    return [
		 	      {
		 		        source: '/api/:path*',
		 		        destination: `http://localhost/api/:path*`,
		 	      },
		    ]
	  },
};

export default nextConfig;
