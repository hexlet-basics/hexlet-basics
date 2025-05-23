class PopulateShouldBeLead < ActiveRecord::Migration[8.0]
  def change
    items = Survey::Item.active.where(slug: [
      "career-change-contact-method-item2",
      "career-change-contact-method-item3",
      "career-change-contact-method-item4"
    ])
    answers = Survey::Answer.where(survey_item: items)
    answers.find_each do |answer|
      user = answer.user
      puts user.first_name
      next if user.leads.any?
      user.tag_list.add('should_be_lead')
      user.save!
    end
  end
end
