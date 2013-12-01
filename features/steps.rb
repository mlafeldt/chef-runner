Given(/^a cookbook named "([^"]*)"$/) do |name|
  step %(a file named "metadata.rb" with:), %(name "#{name}")
  step %(an empty file named "Vagrantfile")
end

Given(/^a cookbook named "([^"]*)" with the recipes? "([^"]*)"$/) do |cookbook, recipes|
  step %(a cookbook named "#{cookbook}")
  recipes.split(",").each do |recipe|
    step %(an empty file named "recipes/#{recipe}.rb")
  end
end

# Insert for debugging
Then(/^shell$/) do
  in_current_dir { system '/bin/bash -i' }
end

Then(/^"([^"]*)" should be run$/) do |cmd|
  history.should include(cmd)
end

Then(/^"([^"]*)" should not be run$/) do |cmd|
  history.all? { |h| h.should_not include(cmd) }
end

Then(/^\/(.*?)\/ should be run$/) do |pattern|
  history.grep(/#{pattern}/).should have(1).item
end

Then(/^\/(.*?)\/ should not be run$/) do |pattern|
  history.all? { |h| h.should_not match /#{pattern}/ }
end

Then(/^"([^"]*)" should be run with the option "([^"]*)"$/) do |cmd, option|
  regexp = /\b#{cmd}\b.*\s#{option}(\s|\z)/
  step %(/#{regexp}/ should be run)
end

Then(/^the runlist should be "([^"]*)"$/) do |runlist|
  step %("chef-solo" should be run with the option "--override-runlist=#{runlist}")
end
