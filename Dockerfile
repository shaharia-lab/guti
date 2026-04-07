FROM ghcr.io/actions/jekyll-build-pages:v1.0.13

# Set workdir to a neutral location
WORKDIR /srv/jekyll-app

# Create separate directories for source and destination
RUN mkdir -p /srv/jekyll-app/source /srv/jekyll-app/dest

# Note: We don't need to COPY docs/ here since we're using a volume mount

# Expose the default Jekyll port
EXPOSE 4000

# Command to build and serve from the correct directories
CMD ["github-pages", "serve", \
     "--source", "/srv/jekyll-app/source", \
     "--destination", "/srv/jekyll-app/dest", \
     "--host", "0.0.0.0", \
     "--force_polling"]