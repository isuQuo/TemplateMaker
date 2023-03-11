# Run the server
`copper run -watch`

# Run the Tailwind server
`cd web`
`npm run dev`

# Create a new route for page
`copper scaffold:route -handler=HandleEditPage -path=/edit app`

# Create a new route for a post page
`copper scaffold:route -handler=HandleSubmitPost -method=Post -path=/submit app`