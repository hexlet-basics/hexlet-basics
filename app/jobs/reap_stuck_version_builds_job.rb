# typed: strict

# Sweeps course version builds that got stuck in `building`. When a build
# worker is hard-killed (OOM/SIGKILL), ExerciseLoader#run never reaches its
# rescue, so the version stays in `building` forever with an empty result and
# no trace anywhere. This periodic reaper (see config/recurring.yml) marks such
# versions `failed` with an explicit result, making the failure visible in the
# admin versions table instead of silent and eternal.
class ReapStuckVersionBuildsJob < ApplicationJob
  extend T::Sig

  sig { void }
  def perform
    Language::Version.stuck_building.find_each do |version|
      version.result = "Build reaped: stuck in 'building' since #{version.updated_at.iso8601} (worker likely killed, e.g. OOM)"
      version.mark_as_failed
      version.save(validate: false)

      Rails.logger.warn("Reaped stuck version build ##{version.id} (#{version.language.slug})")
    end
  end
end
