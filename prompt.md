Fix the frontend JavaScript error:

Cannot read properties of null (reading 'addEventListener')

The app currently fails when I try to enter the demo account or create a new user. This is probably because JavaScript is attaching event listeners to elements that do not exist on the current page, or because IDs in HTML and JS do not match.

Do not rebuild the project.

Tasks:

1. Search all frontend JavaScript files in web/js/ for every use of:
   addEventListener

2. For every addEventListener, make it safe:
   - Get the element first.
   - Check if the element exists before attaching the event listener.
   - Do not allow one missing element to crash the whole app.

Example pattern:

const btn = document.getElementById("buttonId");
if (btn) {
  btn.addEventListener("click", handleClick);
}

3. Wrap page initialization inside DOMContentLoaded where needed:

document.addEventListener("DOMContentLoaded", () => {
  // page setup here
});

4. Check app.html and the JavaScript files for ID mismatches.

Make sure the following buttons/forms exist in HTML and match the JavaScript selectors:
- demo account / enter demo button
- create user form
- create user button
- consent checkbox
- profile save button
- user switch/select controls
- logout/switch user button
- load demo data button
- dashboard tab button
- transactions tab button
- goals tab button
- affordability tab button
- Wizz chat tab button
- data/privacy tab button

5. Fix the demo account flow.

Expected behaviour:
- If no active user exists, show the onboarding/select-user screen.
- Clicking “Enter Demo Account” or “Load Demo User” should create/select a demo user.
- The app should store the active user ID in localStorage.
- After demo user is selected, the dashboard/app sections should become available.
- Refreshing the browser should not throw any JavaScript errors.

6. Fix the create new user flow.

Expected behaviour:
- User enters name/email or username.
- User accepts consent.
- Click create user.
- Backend creates the user or selects existing user.
- Frontend stores active user ID.
- App loads dashboard for that specific user.
- No “Cannot read properties of null” error appears.

7. If a JavaScript file is only meant for app.html, make sure it is not loaded on index.html, privacy.html, terms.html, or ethical-ai.html unless it has null checks.

8. Add console-friendly error messages but do not break the UI.

9. After fixing, test this flow:
   - Open app.html
   - Enter demo account
   - Refresh page
   - Create new user
   - Switch user
   - Add goal
   - Refresh again
   - Confirm there are no console errors

Return:
- Files changed
- What IDs were mismatched
- What was causing the null addEventListener error
- How to test the fix